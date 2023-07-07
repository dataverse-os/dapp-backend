import { createAbstractCompositeDefinition } from "@composedb/devtools";

const schema = JSON.parse(process.argv[2]);
try {
  createAbstractCompositeDefinition(schema);
} catch (error) {
  console.log((error as Error).toString().split("\n")[0]);
}