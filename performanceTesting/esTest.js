const faker = require('faker')
const nanoid = require('nanoid/non-secure').nanoid
const fetch = require('isomorphic-fetch')
const fs = require('fs')
const nread = require('n-readlines')

const uploadItem = async () => {
  const userObj = {
    name: faker.name.findName(),
    email: faker.internet.email(),
    address: faker.address.streetAddress(),
    country: faker.address.country(),
    city: faker.address.city(),
    zipcode: faker.address.zipCode(),
    signup: faker.date.between('2015-01-01', '2020-10-31'),
    bitcoinAddr: faker.finance.bitcoinAddress(),
    uid: nanoid()
  }
  const start = process.hrtime.bigint()
  const resp = await fetch(`http://${process.env.HOSTNAME}:9200/testindex/_doc`, {
    method: 'POST',
    headers: {
      'content-type': 'application/json'
    },
    body: JSON.stringify(userObj)
  })
  if (resp.status > 299) {
    console.log(await resp.text())
  }
  const end = process.hrtime.bigint()
  const diffTime = end - start
  fs.appendFile('./operationTimes.csv', `${Number(diffTime) / 1000000}\n`, () => {}) // milliseconds
  return userObj
}

const uploadTest = async (rounds) => {
  for (let i = 0; i < rounds; i++) {
    const item = await uploadItem()
    // Store every 100 items
    if (i % 100 === 0) {
      console.log(i)
      fs.appendFile('./randomItems.txt', `${JSON.stringify(item)}\n`, () => {})
    }
  }
}

const searchItem = async (fieldName, itemField) => {
  const start = process.hrtime.bigint()
  const resp = await fetch(`http://${process.env.HOSTNAME}:9200/testindex/_search?q=${fieldName}:${itemField}`, {
    method: 'GET'
  })
  if (resp.status > 299) {
    console.log(await resp.text())
  }
  const end = process.hrtime.bigint()
  const diffTime = end - start
  fs.appendFile('./searchItems.csv', `${Number(diffTime) / 1000000}\n`, () => {}) // milliseconds
}

const randomProperty = function (obj) {
  delete obj.signup // ES doesn't like how this is formatted when using GET requests
  const keys = Object.keys(obj)
  return obj[keys[keys.length * Math.random() << 0]]
}

const searchTest = async () => {
  const items = []
  const liner = new nread('./randomItems.txt')
  let line
  while (line = liner.next()) {
    items.push(JSON.parse(line))
  }
  for (let i = 0; i < 1000; i++) {
    // get random item field
    const theKey = randomProperty(items[i])
    await searchItem(theKey, items[i][theKey])
  }
}

const main = async () => {
  // await uploadTest(100000)
  await searchTest()
}

main()
