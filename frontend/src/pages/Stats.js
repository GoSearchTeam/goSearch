import { Badge } from 'react-bootstrap'

const nodes = [
    {
        LocalCluster: "Example Cluster",
        GlobalCluster: null,
        IP: "192.168.1.1",
        Port: 4444,
        Name: "Example Node 1"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: null,
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 2"
    }
]

export default function Stats(props) {
    return (
        <>
            <header className="Section-header">Genaral Status</header>
            <br></br>
            Nodes Running <Badge variant="dark">{getNodes()}</Badge><br></br>
            Documents on Disk <Badge variant="dark">{getDocs()}</Badge><br></br>
            Indexes Currently Searchable <Badge variant="dark">{getIndexes()}</Badge><br></br>
            <hr></hr>
            <header className="Section-header">Node Status</header>
            {getNodeStats()}
        </>
    );
}

function getNodes() {
    return nodes.length;
}

function getDocs() {
    return Math.floor(Math.random() * 1000000) + 100;
}

function getIndexes() {
    return Math.floor(Math.random() * 30) + 1;
}

function getNodeStats(){
    let nodeStats = [];
    nodes.forEach(node => {
        nodeStats.push(
            <div>
                <hr className="short-hr"></hr>
                Name: <Badge variant="dark">{node.Name}</Badge><br></br>
                Local Cluster: <Badge variant="dark">{node.LocalCluster}</Badge><br></br>
                Global Cluster: <Badge variant="dark">{node.GlobalCluster}</Badge><br></br>
                IP: <Badge variant="dark">{node.IP}</Badge><br></br>
                Port: <Badge variant="dark">{node.Port}</Badge><br></br>
            </div>
        )
    });
    return nodeStats;
}