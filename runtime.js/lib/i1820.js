/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-08-2018
 * |
 * | File Name:     i1820.js
 * +===============================================
 */

class I1820 {
  // decode coverts given buffer object into json serializable object
  decode (data) {
    console.error('please implmenet "decode" function')
    return '18.20 is leaving us alone [default]'
  }

  // encode coverts given object into buffer object
  encode (data) {
    console.error('please implmenet "encode" function')
    return Buffer.from('18.20 is leaving us alone [default]', 'ascii')
  }
}

module.exports = I1820
