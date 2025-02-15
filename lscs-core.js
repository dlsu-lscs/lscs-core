import axios from 'axios'

class lscslib {
  constructor(API_KEY) {
    this.API_KEY = null;
  }

  async findMemberByEmail(email) {
    try {
      const response = await axios.post(
        'https://auth.app.dlsu-lscs.org/member',
        { email: email },
        {
          headers: {
            Authorization: `Bearer ${this.API_KEY}`,
            'Content-Type': 'application/json',
          },
        }
      )
      return response.data;
    }
    catch (err) {
      if (err.status === 400) {
        return { error: 'Invalid member' };
      }
      console.error(err)
    }
  }

  async findMemberById(memberid) { // NOTE: DLSU ID!
    try {
      const response = await axios.post(
        'https://auth.app.dlsu-lscs.org/member-id',
        { email: memberid },
        {
          headers: {
            Authorization: `Bearer ${this.API_KEY}`,
            'Content-Type': 'application/json',
          },
        }
      )
      return response.data;
    }
    catch (err) {
      if (err.status === 400) {
        return { error: 'Invalid member' };
      }
      console.error(err)
    }
  }

  async findMemberById(memberid) { // NOTE: DLSU ID!
    try {
      const response = await axios.post(
        'https://auth.app.dlsu-lscs.org/member-id',
        { email: memberid },
        {
          headers: {
            Authorization: `Bearer ${this.API_KEY}`,
            'Content-Type': 'application/json',
          },
        }
      )
      return response.data;
    }
    catch (err) {
      if (err.status === 400) {
        return { error: 'Invalid member' };
      }
      console.error(err)
    }
  }

  async getCommittees() { // NOTE: DLSU ID!
    try {
      const response = await axios.get(
        'https://auth.app.dlsu-lscs.org/committees', {
        headers: {
          Authorization: `Bearer ${this.API_KEY}`,
        },
      })
      return response.data;
    }
    catch (err) {
      if (err.status === 400) {
        console.error(err)
        return { error: 'No access' };
      }
      console.error(err)
    }
  }

}

export default lscslib
