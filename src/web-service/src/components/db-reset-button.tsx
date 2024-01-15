import { resetDatabase } from "@/repository/test-data";
import { Button } from "./ui/button";

export const DBResetButton = () => {
  return (
    <Button
      variant="outline"
      onClick={async () => {
        resetDatabase().then(() => location.reload());
      }}
    >
      Load/Reset Test Data
    </Button>
  );
};
