export namespace algorithms {
	
	export enum AlgorithmType {
	    GeneticAlgorithm = "GA",
	    AHA = "AHA",
	    MOAHA = "MOAHA",
	    GWO = "GWO",
	    oMOAHA = "oMOAHA",
	    MOGWO = "MOGWO",
	    NSGAII = "NSGA-II",
	}

}

export namespace conslay_predetermined {
	
	export class LocFac {
	    locName: string;
	    facName: string;
	
	    static createFrom(source: any = {}) {
	        return new LocFac(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.locName = source["locName"];
	        this.facName = source["facName"];
	    }
	}

}

export namespace data {
	
	export enum ObjectiveType {
	    SafetyObjective = "Safety Objective",
	    HoistingObjective = "Hoisting Objective",
	    RiskObjective = "Risk Objective",
	    TransportCostObjective = "Transport Cost Objective",
	    SafetyHazardObjective = "Safety Hazard Objective",
	    ConstructionCostObjective = "Construction Cost Objective",
	}
	export enum ConstraintType {
	    Overlap = "Overlap",
	    OutOfBound = "OutOfBound",
	    CoverInCraneRadius = "CoverInCraneRadius",
	    InclusiveZone = "InclusiveZone",
	    Size = "Size",
	}
	export enum ProblemName {
	    ContinuousConstructionLayout = "Continuous Construction Layout",
	    GridConstructionLayout = "Grid Construction Layout",
	    PredeterminedConstructionLayout = "Predetermined Construction Layout",
	}

}

export namespace main {
	
	export enum EventType {
	    ProgressEvent = "ProgressEvent",
	    ResultEvent = "ResultEvent",
	}
	export enum CommandType {
	    ExportResult = "ExportResult",
	    SaveChart = "SaveChart",
	}
	export class AlgorithmInput {
	    algorithmName: algorithms.AlgorithmType;
	    algorithmConfig: any;
	
	    static createFrom(source: any = {}) {
	        return new AlgorithmInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.algorithmName = source["algorithmName"];
	        this.algorithmConfig = source["algorithmConfig"];
	    }
	}
	export class ConstraintInput {
	    constraintName: data.ConstraintType;
	    constraintConfig: any;
	
	    static createFrom(source: any = {}) {
	        return new ConstraintInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.constraintName = source["constraintName"];
	        this.constraintConfig = source["constraintConfig"];
	    }
	}
	export class ConstraintsConfigResponse {
	    outOfBoundary?: any;
	    overlap?: any;
	    coverInCraneRadius?: any;
	    inclusiveZone?: any;
	    size?: any;
	
	    static createFrom(source: any = {}) {
	        return new ConstraintsConfigResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outOfBoundary = source["outOfBoundary"];
	        this.overlap = source["overlap"];
	        this.coverInCraneRadius = source["coverInCraneRadius"];
	        this.inclusiveZone = source["inclusiveZone"];
	        this.size = source["size"];
	    }
	}
	export class ObjectiveConfigResponse {
	    risk?: any;
	    hoisting?: any;
	    safety?: any;
	    transportCost?: any;
	    safetyHazard?: any;
	    constructionCost?: any;
	
	    static createFrom(source: any = {}) {
	        return new ObjectiveConfigResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.risk = source["risk"];
	        this.hoisting = source["hoisting"];
	        this.safety = source["safety"];
	        this.transportCost = source["transportCost"];
	        this.safetyHazard = source["safetyHazard"];
	        this.constructionCost = source["constructionCost"];
	    }
	}
	export class ObjectiveInput {
	    objectiveName: data.ObjectiveType;
	    objectiveConfig: any;
	
	    static createFrom(source: any = {}) {
	        return new ObjectiveInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.objectiveName = source["objectiveName"];
	        this.objectiveConfig = source["objectiveConfig"];
	    }
	}
	export class ProblemInput {
	    problemName: data.ProblemName;
	    layoutLength?: number;
	    layoutWidth?: number;
	    facilitiesFilePath?: string;
	    phasesFilePath?: string;
	    gridSize?: number;
	    numberOfLocations?: number;
	    numberOfFacilities?: number;
	    fixedFacilities?: conslay_predetermined.LocFac[];
	
	    static createFrom(source: any = {}) {
	        return new ProblemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.problemName = source["problemName"];
	        this.layoutLength = source["layoutLength"];
	        this.layoutWidth = source["layoutWidth"];
	        this.facilitiesFilePath = source["facilitiesFilePath"];
	        this.phasesFilePath = source["phasesFilePath"];
	        this.gridSize = source["gridSize"];
	        this.numberOfLocations = source["numberOfLocations"];
	        this.numberOfFacilities = source["numberOfFacilities"];
	        this.fixedFacilities = this.convertValues(source["fixedFacilities"], conslay_predetermined.LocFac);
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

