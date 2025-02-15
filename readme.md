# lscs-core
This library enables users to utilize the LSCS Core API in NodeJS hassle free.

# Usage

```js
// Initializing an lscs-core instance
import lscscore from 'lscs-core'
const lscs = new lscscore(process.env.LSCS_AUTH_KEY)

const userData = await lscs.findMemberByEmail('max_chavez@dlsu.edu.ph')
console.log(userData)

const userData2 = await lscs.findMemberById('max_chavez@dlsu.edu.ph')
console.log(userData2)

const committees = await lscs.getCommittees()
console.log(committees)
```
