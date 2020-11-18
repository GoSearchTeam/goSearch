import { Badge } from 'react-bootstrap'

export default function Stats(props) {
    return (
        <>
            <header className="Section-header">Genaral Status</header>
            <br></br>
            Nodes Running <Badge variant="dark">{getNodes()}</Badge><br></br>
            Documents on Disk <Badge variant="dark">{getDocs()}</Badge><br></br>
            Indexes Currently Searchable <Badge variant="dark">{getIndexes()}</Badge><br></br>
        </>
    );
}

function getNodes() {
    return Math.floor(Math.random() * 10) + 1;
}

function getDocs() {
    return Math.floor(Math.random() * 1000000) + 100;
}

function getIndexes() {
    return Math.floor(Math.random() * 30) + 1;
}