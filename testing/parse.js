const fs = require('fs')
const fetch = require('isomorphic-fetch')
const rl = require('readline')
const LineByLineReader = require('line-by-line')
const { format } = require('path')

let lineNum = 0
let sum = BigInt(0)
let reqs = 0
let buff = []

const main = async () => {
  var lr = new LineByLineReader('./rawLogs')

  // lr.on('error', function (err) {
  //   // 'err' contains error object
  // })

  lr.on('line', async function (line) {
    // pause emitting of lines...
    lr.pause()
    const rawTime = line.split(/\s+/)[7]
    let timeMS
    if (rawTime.endsWith('ms')) {
      timeMS = rawTime.split('ms')[0]
    } else {
      timeMS = Number(rawTime.split('s')[0] * 1000)
    }
    fs.appendFileSync('./parsedLogs', `${timeMS}\n`)
    lineNum++
    lr.resume()
  })

  lr.on('end', function () {
    // All lines are read, file is closed now.
    console.log('reqs', reqs)
    console.log('sum', sum)
  })
}

main()
