import type { PageLoad } from './$types';
import {ProblemInfo} from "$lib/wailsjs/go/main/App";

export const load: PageLoad = async ({ params }) => {
  // get locations from problems

  const data = await ProblemInfo()

  console.log(data)

  // get the list of locations if it is not predetermined problem

  return {

  };
};