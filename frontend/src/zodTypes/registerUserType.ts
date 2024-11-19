import { z } from "zod";

export const registerUserType = z.object({
    username: z.string({ message: "Username should be provided" }).min(6, { message: "Minimum length of username is 6" }).max(20, { message: "Maximum length of username is 20" }),
    password: z.string({ message: "Password should be provided" }).min(6, { message: "Minimum length of username is 6" }),
    email: z.string({ message: "Email should be provided" }).email({ message: "Invalid email" }),
})

