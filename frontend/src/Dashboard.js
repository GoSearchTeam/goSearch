//import './App.css';
import { Form, Button, Col } from 'react-bootstrap'

export default function Dashboard(props) {
    return (
        <Form>
            <Form.Group controlId="dashboard.Search">
                <h3 className="Section-header">Search Documents</h3>
                <Form.Row>
                    <Col>
                        <Form.Control placeholder="Fields" className="bg-dark" />
                    </Col>
                    <Col>
                        <Form.Control placeholder="Value" className="bg-dark" />
                    </Col> 
                    <Col>
                        <Button variant="secondary" className="bg-dark" type="submit">
                            Search
                        </Button>
                    </Col>
                </Form.Row>
            </Form.Group>
            <Form.Group controlId="dashboard.AddIndex">
                <h3 className="Section-header">Add Documents</h3>
                <Form.Control as="textarea" className="bg-dark" rows={7} />
                <Button variant="secondary" className="bg-dark" type="submit">
                    Add Document
                </Button>
            </Form.Group>
            <Form.Group controlId="dashboard.UpdateIndex">
                <h3 className="Section-header">Update Document</h3>
                <Form.Control placeholder="Document ID" className="bg-dark" />
                <Form.Control as="textarea" className="bg-dark" rows={6} />
                <Button variant="secondary" className="bg-dark" type="submit">
                    Update Document
                </Button>
            </Form.Group>
            <Form.Group controlId="dashboard.Delete">
                <h3 className="Section-header">Delete Document</h3>
                <Form.Row>
                    <Col>
                        <Form.Control placeholder="Document ID" className="bg-dark" />
                    </Col>
                    <Col>
                        <Button variant="danger" type="submit">
                            Delete
                        </Button>
                    </Col>
                </Form.Row>
            </Form.Group>
        </Form>
    );
}

