export namespace main {
	
	export class PortRule {
	    port: number;
	    protected: boolean;
	    allowedProcess: string;
	    autoHealCmd: string;
	    autoHealDir: string;
	
	    static createFrom(source: any = {}) {
	        return new PortRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.protected = source["protected"];
	        this.allowedProcess = source["allowedProcess"];
	        this.autoHealCmd = source["autoHealCmd"];
	        this.autoHealDir = source["autoHealDir"];
	    }
	}
	export class ProcessDetails {
	    pid: number;
	    name: string;
	    cmdline: string[];
	    envVars: Record<string, string>;
	    cwd: string;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.cmdline = source["cmdline"];
	        this.envVars = source["envVars"];
	        this.cwd = source["cwd"];
	        this.username = source["username"];
	    }
	}
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

