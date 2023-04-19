package io.github.jiashunx.toy.LogSearch.utils;

import com.alibaba.fastjson.JSON;
import io.github.jiashunx.toy.LogSearch.type.CommandType;
import org.junit.Assert;
import org.junit.Test;

import java.util.Arrays;

/**
 * @author jiashunx
 */
public class CommandHelperTest {

    @Test
    public void testPrintHelpInfo() {
        CommandHelper.printHelpInfo();
    }

    @Test
    public void testGetCommandType() {
        Assert.assertEquals(CommandType.T1, CommandHelper.getCommandType("sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066"));
        Assert.assertEquals(CommandType.T2, CommandHelper.getCommandType("sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066"));
        Assert.assertEquals(CommandType.T3, CommandHelper.getCommandType("sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066"));
    }

    @Test
    public void testGetCommandArgs() {
        String[] args0 = CommandHelper.getCommandArgs("sit2 newcore ");
        Assert.assertNull(args0);
        String[] args1 = CommandHelper.getCommandArgs("sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066");
        Assert.assertEquals(5, args1.length);
        Assert.assertEquals("sit2", args1[0]);
        Assert.assertEquals("newcore", args1[1]);
        Assert.assertEquals("2023-03-19", args1[2]);
        Assert.assertEquals("17\\:17\\:17", args1[3]);
        Assert.assertEquals("202304181703222101000066", args1[4]);
        String[] args2 = CommandHelper.getCommandArgs("sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
        Assert.assertEquals(4, args2.length);
        Assert.assertEquals("sit2", args2[0]);
        Assert.assertEquals("newcore", args2[1]);
        Assert.assertEquals("grep", args2[2]);
        Assert.assertEquals("2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066", args2[3]);
        String[] args3 = CommandHelper.getCommandArgs("sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
        Assert.assertEquals(4, args3.length);
        Assert.assertEquals("sit2", args3[0]);
        Assert.assertEquals("newcore", args3[1]);
        Assert.assertEquals("bash", args3[2]);
        Assert.assertEquals("cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066", args3[3]);
    }

}
