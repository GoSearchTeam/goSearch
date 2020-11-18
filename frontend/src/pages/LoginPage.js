import { useHistory, Link } from 'react-router-dom'

export default function LoginPage (props) {
  const history = useHistory()
    return (
    <div className='login-outer'>
      <div className='login-inner'>
        <form>
          <h3 className='Section-header'>Log in</h3>

          <div className='form-group Text-entry-short'>
            <label>Email</label>
            <input type='email' className='form-control bg-dark' placeholder='Enter Email' />
          </div>

          <div className='form-group Text-entry-short'>
            <label>Password</label>
            <input type='password' className='form-control bg-dark' placeholder='Enter password' />
          </div>

          <button
            type='submit'
            variant='primary'
            className='btn btn-dark btn-lg btn-block'
            onClick={() => history.push('/admin/managment')}
          >
            Sign in
          </button>
        </form>
      </div>
    </div>

  )
}

// export default LoginPage
