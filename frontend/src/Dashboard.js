//import './App.css';
import Form from 'react-bootstrap/Form'

export default function Dashboard(props) {
    return (
        <Form>
            <Form.Group controlId="exampleForm.ControlTextarea1">
                <Form.Label>Example textarea</Form.Label>
                <Form.Control as="textarea" rows={3} />
            </Form.Group>
        </Form>
    );
}

