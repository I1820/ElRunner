/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 03-08-2018
 * |
 * | File Name:     index.js
 * +===============================================
 */

const program = require('commander');

program
  .version("1.0.0")
  .option("-t, --target <path>", "module path")
  .option("-j, --job <job>", "job type", /^(decode|encode|rule)$/i)
  .parse(process.argv);

// console.log(program.job)
// console.log(program.target)
