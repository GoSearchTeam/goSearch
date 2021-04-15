import { useState } from 'react'
import { Form, Button, Col, Alert } from 'react-bootstrap'
import { ENDPOINTS } from '../config'

export default function Management(props) {
  const [searchError, setSearchError] = useState(false)
  const [showSearch, setShowSearch] = useState(false)
  const [searchResult, setSearchResult] = useState('Loading...')

  const [showAddError, setShowAddError] = useState(false)
  const [addError, setAddError] = useState('')
  const [showAddResult, setShowAddResult] = useState(false)
  const [addResult, setAddResult] = useState('Loading...')

  return <div>
    {/* Search Form */}
    <Form onSubmit={(e) => {
      e.preventDefault()
      const formData = acquireData(e)
      setShowSearch(false)
      setSearchError(false)
      setSearchResult('Loading...')

      fetch(ENDPOINTS.SEARCH, {
        method: 'POST',
        body: JSON.stringify({
          query: formData.value,
          fields: [formData.index]
        })
      }).then(res => res.json())
      .then(data => {
        setShowSearch(true)
        setSearchResult(JSON.stringify(data, null, 2))
      })
      .catch(() => {
        setSearchError(true)
        setShowSearch(false)
        setSearchResult('Loading...')
      })
    }}>
      <h3 className='Section-header'>Search Documents</h3>
      <Form.Row>
        <Col>
          <Form.Control name='index' type='text' placeholder='Index' className='bg-dark' id='search-index' />
        </Col>
        <Col>
          <Form.Control name='value' type='text' placeholder='Value' className='bg-dark' id='search-value' />
        </Col>
        <Col>
          <Button variant='secondary' className='bg-dark' type='submit'>
            Search
          </Button>
        </Col>
      </Form.Row>
      {searchError &&
        <Alert className='search-error' variant='danger' onClose={() => setSearchError(false)} dismissible>
          An Error occurred
        </Alert>
      }
      {showSearch &&
        <Alert
          className='search-error'
          variant='dark'
          onClose={() => {
            setShowSearch(false)
            setSearchResult('Loading...')
          }}
          dismissible
        >
          <pre>{searchResult}</pre>
        </Alert>
      }
    </Form>
    <br/>

    {/* Add Documents Form */}
    <Form onSubmit={(e) => {
      e.preventDefault()

      setAddError('')
      setShowAddError(false)
      setAddResult('Loading...')
      setShowAddResult(false)

      const { document } = acquireData(e)

      try {
        JSON.parse(document)
      } catch {
        setShowAddError(true)
        setAddError('Please enter Valid JSON.')
        return
      }

      fetch(ENDPOINTS.ADD, {
        method: 'POST',
        body: document
      }).then(res => res.json())
      .then(data => {
        setShowAddResult(true)
        setAddResult(`${data.msg}: ${data.docID}`)
      }).catch(() => {
        setShowAddError(true)
        setAddError('An error occurred when attempting to add the document')
      })
    }}>
      <h3 className='Section-header'>Add Documents</h3>
      <Form.Control name='document' as='textarea' placeholder='Valid JSON' className='bg-dark' rows={7} />
      <Button variant='secondary' className='bg-dark' type='submit'>
        Add Document
      </Button>
      {showAddError &&
        <Alert className='add-error' variant='danger' onClose={() => setShowAddError(false)} dismissible>
          {addError}
        </Alert>
      }
      {showAddResult &&
        <Alert
          className='add-result'
          variant='dark'
          onClose={() => {
            setShowAddResult(false)
            setAddResult('Loading...')
          }}
          dismissible
        >
          {addResult}
        </Alert>

      }
    </Form>
    <br/>

    {/* Update Doucments form */}
    <Form onSubmit={(e) => {
      e.preventDefault()
      const formData = acquireData(e)
      console.log(formData)
    }}>
      <h3 className='Section-header'>Update Document</h3>
      <Form.Control name='docId' type='text' placeholder='Document ID' className='bg-dark' />
      <Form.Control name='docValue' as='textarea' placeholder='Valid JSON' className='bg-dark' rows={6} />
      <Button variant='secondary' className='bg-dark' type='submit'>
        Update Document
      </Button>
    </Form>
    <br/>

    {/* Delete documents form */}
    <Form onSubmit={(e) => {
      e.preventDefault()
      const formData = acquireData(e)
      console.log(formData)
    }}>
      <h3 className='Section-header'>Delete Document</h3>
      <Form.Row>
        <Col>
          <Form.Control name='docId' type='text' placeholder='Document ID' className='bg-dark' />
        </Col>
        <Col>
          <Button variant='danger' type='submit'>
            Delete
          </Button>
        </Col>
      </Form.Row>
    </Form>
    {/* <Form.Group controlId='dashboard.AddIndex'>
          
        </Form.Group>
        <Form.Group controlId='dashboard.UpdateIndex'>
          
        </Form.Group>
        <Form.Group controlId='dashboard.Delete'>
          
        </Form.Group>
      </Form> */}
  </div>
}

function acquireData(event) {
  return Object.fromEntries(new FormData(event.target).entries())
}
