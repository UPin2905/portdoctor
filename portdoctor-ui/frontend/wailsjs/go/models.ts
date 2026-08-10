export namespace port {
	
	export class PortInfo {
	    Port: number;
	    Protocol: string;
	    Address: string;
	    PID: number;
	    Status: string;
	
	    static createFrom(source: any = {}) {
	        return new PortInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Port = source["Port"];
	        this.Protocol = source["Protocol"];
	        this.Address = source["Address"];
	        this.PID = source["PID"];
	        this.Status = source["Status"];
	    }
	}

}

