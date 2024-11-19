import { useForm } from "react-hook-form";
import { Button } from "./ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { registerUserType } from "@/zodTypes/registerUserType";

export const SignUp = () => {
    const navigate = useNavigate();
    const [isSending, setIsSending] = useState<boolean>(false);

    const form = useForm<z.infer<typeof registerUserType>>({
        resolver: zodResolver(registerUserType),
        defaultValues: {
            username: "",
            email: "",
            password: "",
        },
    });
    const { toast } = useToast();

    async function onSubmit(values: z.infer<typeof registerUserType>) {
        try {
            setIsSending(true);
            const res = await axios.post(
                "http://localhost:8000/api/v1/user/sign-up",
                values
            );
            console.log(res);
            if (res.status != 200) {
                toast({
                    title: "Issue signing in",
                    description: `${res.data.message}`,
                });
                return;
            }
            toast({
                title: "Signup successful",
            });
            navigate("/signin");
            // redirect to singin page
        } catch (err: any) {
            toast({
                title: "Issue signing in",
                description: `There was an error signing up ${err.message}`,
                variant: "destructive"
            });
        } finally {
            setIsSending(false);
        }
    }

    return (
        <div className="flex items-center justify-center h-screen">
            <div className="md:border-2 p-4">
                <h1 className="font-bold text-2xl pb-4 text-center">Sign-up to use Gather Town</h1>
                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit(onSubmit)}
                        className="space-y-8"
                    >
                        <div className="w-60 md:w-96">
                            <FormField
                                control={form.control}
                                name="username"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Username</FormLabel>
                                        <FormControl>
                                            <Input
                                                className="h-12 md:text-xl"
                                                placeholder="username"
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <div className="w-60 md:w-96">
                            <FormField
                                control={form.control}
                                name="email"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Email</FormLabel>
                                        <FormControl>
                                            <Input
                                                className="h-12 md:text-xl"
                                                placeholder="email"
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <div className="w-60 md:w-96">
                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Password</FormLabel>
                                        <FormControl>
                                            <Input
                                                type="password"
                                                placeholder="password"
                                                {...field}
                                                className="h-12 md:text-xl"
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        {!isSending && <Button type="submit">Register</Button>}
                    </form>
                </Form>
                <p className="pt-4">Already have an account? <p className="text-blue-500 cursor-pointer" onClick={() => { navigate("/signin") }}>Signin</p></p>
            </div>
        </div>
    );
};
