export namespace appconfig {
	
	export class ChatState {
	    selectedApiIndex: number;
	    selectedModel: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selectedApiIndex = source["selectedApiIndex"];
	        this.selectedModel = source["selectedModel"];
	    }
	}
	export class KeyRow {
	    protocol: string;
	    alias: string;
	    apiKey: string;
	    baseUrl: string;
	    success: boolean;
	    models: string[];
	    modelCount: number;
	    errorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new KeyRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.alias = source["alias"];
	        this.apiKey = source["apiKey"];
	        this.baseUrl = source["baseUrl"];
	        this.success = source["success"];
	        this.models = source["models"];
	        this.modelCount = source["modelCount"];
	        this.errorMsg = source["errorMsg"];
	    }
	}
	export class AppConfig {
	    keysList: KeyRow[];
	    chatState: ChatState;
	    defaults: Record<string, string>;
	    startupPass: string;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keysList = this.convertValues(source["keysList"], KeyRow);
	        this.chatState = this.convertValues(source["chatState"], ChatState);
	        this.defaults = source["defaults"];
	        this.startupPass = source["startupPass"];
	        this.prompt = source["prompt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class ModelCache {
	    apiModels: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ModelCache(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiModels = source["apiModels"];
	    }
	}

}

export namespace checker {
	
	export class CheckResult {
	    alias: string;
	    key: string;
	    protocol: string;
	    isValid: boolean;
	    models: string[];
	    errorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.alias = source["alias"];
	        this.key = source["key"];
	        this.protocol = source["protocol"];
	        this.isValid = source["isValid"];
	        this.models = source["models"];
	        this.errorMsg = source["errorMsg"];
	    }
	}
	export class BatchCheckResult {
	    results: CheckResult[];
	
	    static createFrom(source: any = {}) {
	        return new BatchCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.results = this.convertValues(source["results"], CheckResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

