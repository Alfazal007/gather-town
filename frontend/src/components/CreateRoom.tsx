import { useForm } from "react-hook-form";
import { Button } from "./ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { createRoomUserType } from "@/zodTypes/createRoomType";
import { UserContext } from "@/context/UserContext";

export const Createroom = () => {
    const navigate = useNavigate();
    const [isSending, setIsSending] = useState<boolean>(false);
    const { user } = useContext(UserContext)

    const form = useForm<z.infer<typeof createRoomUserType>>({
        resolver: zodResolver(createRoomUserType),
        defaultValues: {
            name: "",
        },
    });
    const { toast } = useToast();

    async function onSubmit(values: z.infer<typeof createRoomUserType>) {
        try {
            setIsSending(true);
            console.log({ values })
            const res = await axios.post(
                "http://localhost:8000/api/v1/room/create-room",
                values,
                {
                    headers: {
                        Authorization: `Bearer ${user?.accessToken}`
                    }
                }
            );
            console.log(res);
            if (res.status != 201) {
                toast({
                    title: "Issue creating the room",
                    description: `${res.data.message}`,
                    variant: "destructive"
                });
                return;
            }
            toast({
                title: "Room created successfully",
            });
            navigate("/");
        } catch (err: any) {
            console.log(err)
            toast({
                title: "Issue creating the room",
                description: `There was an error creating the room ${err.message}`,
                variant: "destructive"
            });
        } finally {
            setIsSending(false);
        }
    }

    useEffect(() => {
        if (!user) {
            navigate("/signin")
            return
        }
    }, [user])

    return (
        < div className="flex items-center justify-center h-screen" >
            <div className="md:border-2 p-4">
                <h1 className="font-bold text-2xl pb-4 text-center">Give a name to this room</h1>
                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit(onSubmit)}
                        className="space-y-8"
                    >
                        <div className="w-60 md:w-96">
                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormControl>
                                            <Input
                                                className="h-12 md:text-xl"
                                                placeholder="room name"
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        {!isSending && <Button type="submit">Create</Button>}
                    </form>
                </Form>
            </div>
        </div >
    );
};
