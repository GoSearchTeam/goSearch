import React from 'react';
import { Badge } from 'react-bootstrap'

const dummyNodes = [
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

export class Stats extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            nodeCount: getNodes(),
            docCount: getDocs(),
            indexCount: getIndexes()
        }
    }

    render() {
        return (
            <>
                <div className="body">
                    <header className="Section-header">General Status</header>
                    <br></br>
                    Nodes Running <Badge variant="dark">{this.state.nodeCount}</Badge><br></br>
                    Documents on Disk <Badge variant="dark">{this.state.docCount}</Badge><br></br>
                    Indexes Currently Searchable <Badge variant="dark">{this.state.indexCount}</Badge><br></br>
                    <hr className="grey-hr"></hr>
                    <header className="Section-header">Node Status</header>
                    {getNodeStats()}
                </div>
            </>
        );
    }
    
    componentDidMount() {
        this.interval = setInterval(() => {
            this.setState({
                nodeCount: getNodes(),
                docCount: getDocs(),
                indexCount: getIndexes()
            })
        }, 5e3);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

}

function getNodes() {
    return dummyNodes.length;
}

function getDocs() {
    return Math.floor(Math.random() * 1000000) + 100;
}

function getIndexes() {
    return Math.floor(Math.random() * 30) + 1;
}

function getNodeStats() {
    let nodeStats = [];
    dummyNodes.forEach(node => {
        nodeStats.push(
            <div className="body">
                <hr className="short-hr grey-hr"></hr>
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