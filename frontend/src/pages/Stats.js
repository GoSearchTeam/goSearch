import React from 'react';
import { Badge } from 'react-bootstrap'
import NodeStats from '../components/NodeStats';
import { ENDPOINTS } from '../config';

export class Stats extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            nodes: [],
            itemCount: 0,
            indexCount: 0
        }
    }

    render() {
        return (
            <>
                <div className="body">
                    <header className="Section-header">General Status</header>
                    <br></br>
                    <div className="general-stats-gird">
                        <div className="general-stats">
                            Nodes Running <Badge variant="dark" className="general-stats-badge">{this.state.nodes.length}</Badge>
                        </div>
                        <div className="general-stats">
                            Items available <Badge variant="dark" className="general-stats-badge">{this.state.itemCount}</Badge>
                        </div>
                        <div className="general-stats">
                            Indexes Currently Searchable <Badge variant="dark" className="general-stats-badge">{this.state.indexCount}</Badge>
                        </div>
                    </div>
                    <hr className="grey-hr"></hr>
                    <header className="Section-header">Node Status</header>
                    <div className="node-stats-grid">
                        {this.state.nodes.map(node => <NodeStats node = {node}/>)}
                    </div>
                </div>
            </>
        );
    }
    
    async componentDidMount() {
        this.setState({
            nodes: await getNodes(),
            itemCount: await getItems(),
            indexCount: await getIndexes()
        });

        this.interval = setInterval(async() => {
            this.setState({
                nodes: await getNodes(),
                itemCount: await getItems(),
                indexCount: await getIndexes()
            });
        }, 5e3);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

}

async function getNodes() {
    return await fetch(ENDPOINTS.NODES)
        .then(res => res.json())
}

async function getItems() {
    // return await fetch(ENDPOINTS.LIST_ITEMS)
    //     .then(res => res.json())
    //     .then(data => {
    //         let count = 0;
    //         if (data)
    //             data.forEach(item => count += item.IndexValues.length)
    //         return count;
    //      });
    return 1;
}

async function getIndexes() {
    return await fetch(ENDPOINTS.LIST_IDXS)
        .then(res => res.json())
        .then(data => data?.length || 0);
}
