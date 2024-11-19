import { z } from "zod";

export const loginUserType = z.object({
    username: z.string({ message: "Username or email should be provided" }),
    password: z.string({ message: "Password should be provided" }).min(6, { message: "Minimum length of username is 6" }),
})
