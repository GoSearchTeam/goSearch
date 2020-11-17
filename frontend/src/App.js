import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css'
import { BrowserRouter as Router, Switch, Route, useHistory } from 'react-router-dom'
import LoginPage from './pages/LoginPage.js'
import Dashboard from './pages/Dashboard.js'

function App () {
  return (
    <Router>
      <div className='App'>
        <header className='App-header'>
          Go Search (Admin Land)
        </header>
        <Switch>
          {/* THE ORDER HERE MATTERS, which ever matches first will load i.e. /admin will load /admin/dashboard if put first */}
          <Route path='/admin/dashboard' render={props => <Dashboard {...props} />} />
          <Route path='/admin' render={props => <LoginPage {...props} />} />
        </Switch>
      </div>
    </Router>
  )
}

export default App
