
import hello.Hello;

public class Main {
    public static void main(String[] args) {
        System.loadLibrary("gojni");
        Hello.test2(new byte[]{-1, 0, 1, 2, 3});
        System.out.println("Hello.getMessage() = " + Hello.getMessage());
    }
}
