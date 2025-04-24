import type {PageLoad} from "../../../.svelte-kit/types/src/routes/data/$types";
import {ObjectivesInfo, ProblemInfo} from "$lib/wailsjs/go/main/App";

export const load: PageLoad = async ({ params }) => {
  // get locations from problems

  const data = await ProblemInfo()
  // get the list of locations if it is not predetermined problem

  return {
    problemInfo: data
  };
};