import axios from "axios";
import toast from "react-hot-toast";

export const resetDatabase = async () => {
  try {
    await axios.post("/api/v1/reset");
  } catch (error) {
    toast.error(`Error resetting database.`);
    console.log(error);
  }
};
