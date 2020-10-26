const csv = require('csv-parser')
const fs = require('fs')
const fetch = require('isomorphic-fetch')
const rl = require('readline')
const LineByLineReader = require('line-by-line')

let lineNum = 0
let sum = BigInt(0)
let reqs = 0
let buff = []

const main = async () => {
  var lr = new LineByLineReader('/Users/dangoodman/Downloads/geographic-units-by-industry-and-statistical-area-2000-19-descending-order/Data7602DescendingYearOrder.csv')

  // lr.on('error', function (err) {
  //   // 'err' contains error object
  // })

  lr.on('line', async function (line) {
    // pause emitting of lines...
    lr.pause()
    const lineItem = line.split(',')
    buff.push({
      anzsic06: lineItem[0],
      Area: lineItem[1],
      year: lineItem[2],
      geo_count: lineItem[3],
      ec_count: lineItem[4],
      lineNum: lineNum
    })
    // ...do your asynchronous line processing..
    lineNum++
    if (lineNum >= 50 && lineNum % 20 === 0) {
      const start = process.hrtime.bigint()
      const resp = await fetch('http://localhost:8080/index/addMultiple', {
        method: 'post',
        body: JSON.stringify({
          items: buff
        }),
        headers: {
          'content-type': 'application/json'
        }
      })
      const end = process.hrtime.bigint()
      sum += (end - start)
      reqs++
      console.log(lineNum)
      if (resp.status >= 300) {
        console.log('got status', resp.status)
        process.exit()
      }
      buff = []
    }
    lr.resume()
  })

  lr.on('end', function () {
    // All lines are read, file is closed now.
    console.log('reqs', reqs)
    console.log('sum', sum)
  })
}

main()
