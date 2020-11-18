import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'
import { Navbar, Nav } from 'react-bootstrap'
import LoginPage from './pages/LoginPage.js'
import Managment from './pages/Managment.js'
import Stats from './pages/Stats.js'

function App() {
  return (
    <Router>
      <div className='App'>
        <Navbar bg="dark" variant="dark" expand="lg">
          <Navbar.Brand href="./managment">Go Search</Navbar.Brand>
          <Nav className="mr-auto">
            <Nav.Link href="./managment">Managment</Nav.Link>
            <Nav.Link href="./stats">Status</Nav.Link>
          </Nav>
        </Navbar>
        <br></br>
        <Switch>
          {/* THE ORDER HERE MATTERS, which ever matches first will load i.e. /admin will load /admin/dashboard if put first */}
          <Route path='/admin/managment' render={props => <Managment {...props} />} />
          <Route path='/admin/stats' render={props => <Stats {...props} />} />
          <Route path='/admin' render={props => <LoginPage {...props} />} />
        </Switch>
      </div>
    </Router>
  )
}

export default App
