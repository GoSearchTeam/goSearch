import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import { BrowserRouter as Router, Switch, Route, useHistory } from "react-router-dom";
import LoginPage from "./LoginPage.js";
import Dashboard from "./Dashboard.js";

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          Go Search (Admin Land)
        </header>
        <Switch>
          <Route path="/admin">
            <LoginPage></LoginPage>
          </Route>
          <Route path="/admin/dashboard">
            <Dashboard></Dashboard>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
