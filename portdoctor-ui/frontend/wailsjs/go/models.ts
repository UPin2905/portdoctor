export namespace main {
	
	export class UIPortInfo {
	    port: number;
	    status: string;
	    pid: number;
	    processName: string;
	    project: string;
	    cpu: number;
	    ram: number;
	    sharedUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UIPortInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.status = source["status"];
	        this.pid = source["pid"];
	        this.processName = source["processName"];
	        this.project = source["project"];
	        this.cpu = source["cpu"];
	        this.ram = source["ram"];
	        this.sharedUrl = source["sharedUrl"];
	    }
	}

}

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

