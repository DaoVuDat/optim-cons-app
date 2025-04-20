export namespace algorithms {
	
	export enum AlgorithmType {
	    GeneticAlgorithm = "GA",
	    AHA = "AHA",
	    MOAHA = "MOAHA",
	    GWO = "GWO",
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
	
	    static createFrom(source: any = {}) {
	        return new ConstraintsConfigResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outOfBoundary = source["outOfBoundary"];
	        this.overlap = source["overlap"];
	        this.coverInCraneRadius = source["coverInCraneRadius"];
	        this.inclusiveZone = source["inclusiveZone"];
	    }
	}
	export class ObjectiveConfigResponse {
	    risk?: any;
	    hoisting?: any;
	    safety?: any;
	    transportCost?: any;
	    safetyHazard?: any;
	
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
	    predeterminedLoc?: string;
	
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
	        this.predeterminedLoc = source["predeterminedLoc"];
	    }
	}

}

