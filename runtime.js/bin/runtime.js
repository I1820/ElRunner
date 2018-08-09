/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 03-08-2018
 * |
 * | File Name:     runtime.js
 * +===============================================
 */

const program = require('commander')
const I1820 = require('../lib/i1820')
const readlineSync = require('readline-sync')

program
  .version('1.0.0')
  .option('-t, --target <path>', 'module path')
  .option('-j, --job <job>', 'job type', /^(decode|encode|rule)$/i)
  .parse(process.argv)

// console.log(program.job)
const UI1820 = require(program.target)

let ui1820 = new UI1820()

if (ui1820 instanceof I1820) {
  let line = readlineSync.question()
  let i = Buffer.from(line, 'base64')
  let o = ui1820.decode(i)
  console.log(JSON.stringify(o))
} else {
  console.log('Please corrects your script')
}
