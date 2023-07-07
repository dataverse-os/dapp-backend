import { CeramicClient } from "@ceramicnetwork/http-client";
import { DID } from "dids";
import { Ed25519Provider } from "key-did-provider-ed25519";
import { getResolver } from "key-did-resolver";
import { fromString } from "uint8arrays/from-string";

(async () => {
  const input = JSON.parse(process.argv[2]);
  const privateKey = fromString(input.key!, "base16");
  const did = new DID({
    resolver: getResolver(),
    provider: new Ed25519Provider(privateKey),
  });
  await did.authenticate();
  const ceramic = new CeramicClient(input.ceramic!);
  ceramic.did = did;
  try {
    console.log(await ceramic.admin.nodeStatus())
  } catch (error) {
    console.log((error as Error).toString().split("\n")[0]);
  }
})();