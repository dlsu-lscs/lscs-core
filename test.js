import lscslib from './index.js'

const userData = await lscs.findMemberByEmail('max_chavez@dlsu.edu.ph')
console.log(userData)

const userData2 = await lscs.findMemberById('max_chavez@dlsu.edu.ph')
console.log(userData2)

const committees = await lscs.getCommittees()
console.log(committees)

