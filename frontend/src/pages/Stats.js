import React from 'react';
import { Badge } from 'react-bootstrap'
import NodeStats from '../components/NodeStats';

const dummyNodes = [
    {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "192.168.1.1",
        Port: 4444,
        Name: "Example Node 1"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 2"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 3"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 4"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 5"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 6"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 7"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 8"
    }, {
        LocalCluster: "Example Cluster",
        GlobalCluster: "AWS1",
        IP: "255.255.255.255",
        Port: 5555,
        Name: "Example Node 9"
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

    render(props) {
        return (
            <>
                <div className="body">
                    <header className="Section-header">General Status</header>
                    <br></br>
                    <div className="general-stats-gird">
                        <div className="general-stats">
                            Nodes Running <Badge variant="dark" className="general-stats-badge">{this.state.nodeCount}</Badge>
                        </div>
                        <div className="general-stats">
                            Documents on Disk <Badge variant="dark" className="general-stats-badge">{this.state.docCount}</Badge>
                        </div>
                        <div className="general-stats">
                            Indexes Currently Searchable <Badge variant="dark" className="general-stats-badge">{this.state.indexCount}</Badge>
                        </div>
                    </div>
                    <hr className="grey-hr"></hr>
                    <header className="Section-header">Node Status</header>
                    <div className="node-stats-grid">
                        {dummyNodes.map(node => <NodeStats node = {node}/>)}
                    </div>
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
