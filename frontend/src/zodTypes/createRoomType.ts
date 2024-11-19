import { z } from "zod";

export const createRoomUserType = z.object({
    name: z.string({ message: "Room name should be provided" }).min(4, { message: "Minimum length of roomname is 4" }).max(20, { message: "Maximum length of roomname is 20" }),
})
