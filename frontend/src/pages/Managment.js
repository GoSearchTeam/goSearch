import React from 'react';
import { Form, Button, Col, Alert } from 'react-bootstrap'
import { ENDPOINTS } from '../config'

export class Management extends React.Component {

  constructor(props) {
    super(props)

    this.state = {
      showSearchError: false,
      showSearchResult: false,
      searchResult: 'Loading...',

      showAddError: false,
      addError: 'Error',
      showAddResult: false,
      addResult: 'Loading...',

      showUpdateResult: false,
      updateResult: 'Loading...',
      showUpdateError: false,
      updateError: 'Error',

      showDeleteSuccess: false,
      showDeleteError: false
    }
  }

  render() {
    return (
      <div>
        {/* Search Form */}
        <Form onSubmit={(e) => {
          const formData = acquireData(e)
          this.setState({
            showSearchResult: false,
            showSearchError: false,
            searchResult: 'Loading...'
          })

          fetch(ENDPOINTS.SEARCH, {
            method: 'POST',
            body: JSON.stringify({
              query: formData.value,
              fields: [formData.index]
            })
          }).then(res => res.text())
            .then(text => JSON.parse(text.replace(/("[^"]*"\s*:\s*)(\d{16,})/g, '$1"$2"')))
            .then(data => {
              this.setState({
                showSearchResult: true,
                searchResult: JSON.stringify(data, null, 2)
              })
            })
            .catch((e) => {
              console.error(e)
              this.setState({
                showSearchError: true,
                showSearchResult: false,
                searchResult: 'Loading...'
              })
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
          {this.state.showSearchError &&
            <Alert className='search-error' variant='danger' onClose={() => this.setState({showSearchError: false})} dismissible>
              An Error occurred
            </Alert>
          }
          {this.state.showSearchResult &&
            <Alert
              className='search-error'
              variant='dark'
              onClose={() => {
                this.setState({
                  showSearchResult: false,
                  searchResult: 'Loading...'
                })
              }}
              dismissible
            >
              <pre>{this.state.searchResult}</pre>
            </Alert>
          }
        </Form>
        <br />

        {/* Add Documents Form */}
        <Form onSubmit={(e) => {
          this.setState({
            addError: 'Error',
            showAddError: false,
            addResult: 'Loading',
            showAddResult: false
          })

          const { document } = acquireData(e)

          try {
            JSON.parse(document)
          } catch {
            this.setState({
              showAddError: true,
              addError: 'Please enter Valid JSON.'
            })
            return
          }

          fetch(ENDPOINTS.ADD, {
            method: 'POST',
            body: document
          }).then(res => res.text())
            .then(text => JSON.parse(text.replace(/("[^"]*"\s*:\s*)(\d{16,})/g, '$1"$2"')))
            .then(data => {
              this.setState({
                showAddResult: true,
                addResult: `${data.msg}: ${data.docID}`
              })
            }).catch(() => {
              this.setState({
                showAddError: true,
                addError: 'An error occurred when attempting to add the document'
              })
            })
        }}>
          <h3 className='Section-header'>Add Documents</h3>
          <Form.Control name='document' as='textarea' placeholder='Valid JSON' className='bg-dark' rows={7} />
          <Button variant='secondary' className='bg-dark' type='submit'>
            Add Document
      </Button>
          {this.state.showAddError &&
            <Alert className='add-error' variant='danger' onClose={() => this.setState({showAddError: false})} dismissible>
              {this.state.addError}
            </Alert>
          }
          {this.state.showAddResult &&
            <Alert
              className='add-result'
              variant='dark'
              onClose={() => {
                this.setState({
                  showAddResult: false,
                  addResult: 'Loading...'
                })
              }}
              dismissible
            >
              {this.state.addResult}
            </Alert>

          }
        </Form>
        <br />

        {/* Update Documents form */}
        <Form onSubmit={(e) => {
          this.setState({
            showUpdateResult: false,
            showUpdateError: false,
            updateResult: 'Loading...',
            updateError: 'Error'
          })

          const data = acquireData(e)

          if (!parseInt(data.docId)) {
            this.setState({
              showUpdateError: true,
              updateError: 'Please enter a numeric document ID.'
            })
            return
          }

          try {
            JSON.parse(data.document)
          } catch {
            this.setState({
              showUpdateError: true,
              updateError: 'Please enter Valid JSON.'
            })
            return
          }

          fetch(ENDPOINTS.UPDATE, {
            method: 'POST',
            body: `{"docID":${data.docId}, ${data.document.substring(1)}`
          }).then(res => res.text())
          .then(msg => {
            console.log('message', msg)
            this.setState({
              showUpdateResult: true,
              updateResult: msg
            })
          }).catch(() => {
            this.setState({
              showUpdateError: true,
              updateError: 'An error occurred when trying to update the document.'
            })
          })

        }}>
          <h3 className='Section-header'>Update Document</h3>
          <Form.Control name='docId' type='text' placeholder='Document ID' className='bg-dark' />
          <Form.Control name='document' as='textarea' placeholder='Valid JSON' className='bg-dark' rows={6} />
          <Button variant='secondary' className='bg-dark' type='submit'>
            Update Document
          </Button>
          {this.state.showUpdateError &&
            <Alert className='add-error' variant='danger' onClose={() => this.setState({showUpdateError: false})} dismissible>
              {this.state.updateError}
            </Alert>
          }
          {this.state.showUpdateResult &&
            <Alert className='add-error' variant='dark' onClose={() => this.setState({showUpdateResult: false})} dismissible>
              {this.state.updateResult}
            </Alert>
          }
        </Form>
        <br />

        {/* Delete documents form */}
        <Form onSubmit={(e) => {
          this.setState({
            showDeleteError: false,
            showDeleteSuccess: false
          })
          e.preventDefault()
          const { docId: idString } = acquireData(e)
          fetch(ENDPOINTS.DELETE, {
            method: 'POST',
            body: `{"docID": ${idString}}`
          }).then(res => {
            if (res.ok)
              this.setState({showDeleteSuccess: true})
            else
              this.setState({showDeleteError: true})
          })
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
          {this.state.showDeleteError &&
            <Alert className='search-error' variant='danger' onClose={() => this.setState({showDeleteError: false})} dismissible>
              An Error occurred
            </Alert>
          }
          {this.state.showDeleteSuccess &&
            <Alert className='search-error' variant='dark' onClose={() => this.setState({showDeleteSuccess: false})} dismissible>
              Document Deleted
            </Alert>
          }
        </Form>
      </div>
    )
  }

}


function acquireData(event) {
  event.preventDefault()
  return Object.fromEntries(new FormData(event.target).entries())
}
