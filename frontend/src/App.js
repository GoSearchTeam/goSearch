import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'
import { Navbar, Nav } from 'react-bootstrap'
import LoginPage from './pages/LoginPage.js'
import Management from './pages/Managment.js'
import { Stats } from './pages/Stats.js'

function App() {
  return (
    <Router>
      <div className='App'>
        <Navbar bg="dark" variant="dark" expand="lg" className="nav-bar-round">
          <Navbar.Brand href="./management">Go Search</Navbar.Brand>
          <Nav className="mr-auto">
            <Nav.Link href="./management">Management</Nav.Link>
            <Nav.Link href="./stats">Status</Nav.Link>
          </Nav>
        </Navbar>
        <br></br>
        <Switch>
          {/* THE ORDER HERE MATTERS, which ever matches first will load i.e. /admin will load /admin/dashboard if put first */}
          <Route path='/admin/management' render={props => <Management {...props} />} />
          <Route path='/admin/stats' render={props => <Stats {...props} />} />
          <Route path='/admin' render={props => <LoginPage {...props} />} />
        </Switch>
      </div>
    </Router>
  )
}

export default App
