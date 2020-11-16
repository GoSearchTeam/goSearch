const faker = require('faker')
const nanoid = require('nanoid/non-secure').nanoid
const fetch = require('isomorphic-fetch')
const fs = require('fs')

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
  const resp = await fetch(`http://${process.env.HOSTNAME}:8080/index/add`, {
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
  fs.appendFile('./operationTimes.csv', `ADD-SINGLE, ${Number(diffTime) / 1000000}\n`, () => {}) // milliseconds
  return userObj
}

const uploadTest = async (rounds) => {
  for (let i = 0; i < rounds; i++) {
    const item = await uploadItem()
    // Store every 100 items
    if (i % 100 === 0) {
      fs.appendFile('./randomItems.txt', `${JSON.stringify(item)}\n`, () => {})
    }
  }
}

const main = async () => {
  uploadTest(10)
}

main()
