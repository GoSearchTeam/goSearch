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
                Local Cluster: <Badge variant="dark">{node.LocalCluster || 'none'}</Badge><br></br>
                Global Cluster: <Badge variant="dark">{node.GlobalCluster || 'none'}</Badge><br></br>
                IP: <Badge variant="dark">{node.IP || 'unknown'}</Badge><br></br>
                Port: <Badge variant="dark">{node.Port}</Badge><br></br>
                API Port: <Badge variant="dark">{node.APIPort}</Badge><br></br>
            </div>
        );
    }

}

export default NodeStats;
