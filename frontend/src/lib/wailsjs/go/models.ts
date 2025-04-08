export namespace main {
	
	export class ObjectiveConfigResponse {
	    risk?: any;
	    hoisting?: any;
	    safety?: any;
	
	    static createFrom(source: any = {}) {
	        return new ObjectiveConfigResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.risk = source["risk"];
	        this.hoisting = source["hoisting"];
	        this.safety = source["safety"];
	    }
	}
	export class ObjectiveInput {
	    objectiveName: objectives.ObjectiveType;
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
	    problemName: objectives.ProblemType;
	    layoutLength?: number;
	    layoutWidth?: number;
	    facilitiesFilePath?: string;
	    phasesFilePath?: string;
	    gridSize?: string;
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

export namespace objectives {
	
	export enum ProblemType {
	    ContinuousConstructionLayout = "Continuous Construction Layout",
	    GridConstructionLayout = "Grid Construction Layout",
	    PredeterminedConstructionLayout = "Predetermined Construction Layout",
	}
	export enum ObjectiveType {
	    SafetyObjective = "Safety Objective",
	    HoistingObjective = "Hoisting Objective",
	    RiskObjective = "Risk Objective",
	}
	export enum ConstraintType {
	    Overlap = "Overlap",
	    OutOfBound = "OutOfBound",
	    CoverInCraneRadius = "CoverInCraneRadius",
	    InclusiveZone = "InclusiveZone",
	}

}

