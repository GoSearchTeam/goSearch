import React from 'react';
import { Badge } from 'react-bootstrap';

class NodeStats extends React.Component {

    render() {
        const { node } = this.props;
        return (
            <div className="node-stats">
                <h5>
                    {node.Name}
                </h5>
                <hr className="stats-hr"></hr>
                Local Cluster: <Badge variant="dark">{node.LocalCluster}</Badge><br></br>
                Global Cluster: <Badge variant="dark">{node.GlobalCluster}</Badge><br></br>
                IP: <Badge variant="dark">{node.IP}</Badge><br></br>
                Port: <Badge variant="dark">{node.Port}</Badge><br></br>
            </div>
        );
    }

}

export default NodeStats;
