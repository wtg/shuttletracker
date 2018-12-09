abstract class Exp {
  def visit[A](proc_And: (A,A) => A, proc_Or: (A,A) => A, proc_Not: (A) => A, proc_Const: (Boolean) => A) : A 
}

case class And(val left: Exp, val right: Exp) extends Exp {
  override def visit[A](proc_And: (A,A) => A, proc_Or: (A,A) => A, proc_Not: (A) => A, proc_Const: (Boolean) => A) : A =
    proc_And(left.visit(proc_And,proc_Or,proc_Not,proc_Const),
             right.visit(proc_And,proc_Or,proc_Not,proc_Const))
}

case class Or(val left: Exp, val right: Exp) extends Exp {

  // Scala requires an "override" annotation when overriding a non-abstract method.
  // When implementing an abstract method, we may omit "override".
  def visit[A](proc_And: (A,A) => A, proc_Or: (A,A) => A, proc_Not: (A) => A, proc_Const: (Boolean) => A) : A =
    proc_Or(left.visit(proc_And,proc_Or,proc_Not,proc_Const),
	    right.visit(proc_And,proc_Or,proc_Not,proc_Const))
}
   
case class Not(val op: Exp) extends Exp {
  def visit[A](proc_And: (A,A) => A, proc_Or: (A,A) => A, proc_Not: (A) => A, proc_Const: (Boolean) => A) : A =
    proc_Not(op.visit(proc_And,proc_Or,proc_Not,proc_Const))
}

case class Const(val v: Boolean) extends Exp {
  def visit[A](proc_And: (A,A) => A, proc_Or: (A,A) => A, proc_Not: (A) => A, proc_Const: (Boolean) => A) : A =
    proc_Const(v)
}

object Expressions {

  def main(args : Array[String]) {

    /* YOUR CODE HERE */
    /** eval returns the boolean value of the expression. 
      *  
      * Fill in the body of eval using pattern matching, as we did in class. 
      * It currently returns a dummy value in order to compile.
      */
    
    def eval(exp : Exp) : Boolean = exp match { 
       case And(l,r) => true
       // ...

    }

   
    /* YOUR CODE HERE */
    /**
      *  eval_X functions instantiate proc_X ones. Type argument Boolean
      *  replaces, i.e., instantiates, type parameter A. In this case we 
      *  make use of the higher-order visit method to traverse and evaluate 
      *  the underlying expression.
       
      *  Fill in the code for eval_Or, eval_And, eval_Not, and eval_Const.
      *  All currently return dummy values in order to compile. 
      */
  
    def eval_Or(l:Boolean,r:Boolean) : Boolean = true // ...
    def eval_And(l:Boolean,r:Boolean) : Boolean = true // ...
    def eval_Not(o:Boolean) : Boolean = true // ...
    def eval_Const(v:Boolean) : Boolean = true // ...
    


    def parens(s:String) : String = "("+s+")";

    /* YOUR CODE HERE */
    /** 
      *  infix_X functions instantiate proc_X ones. Type argument (Int,String)
      *  replaces type parameter A. A (Int,String) tuple carries the precedence
      *  (Int-element of the tuple), and the infix string representation 
      *  (String-element of the tuple) of the expression. Precedence ranges from 
      *  0 to 3: 0 if Or-epxression, 1 if And-expression, 2 if Not, and 3 if Const.
      * 
      *  infix_Or is included. Fill in infix_And, infix_Not, and infix_Const.
      *  Again, all currently return dummy values.
     */

    
    def infix_Or(l:(Int,String),r:(Int,String)) : (Int,String) = {
      val (_,op1) = l;
      val (p2,s2) = r;

      // If precedence of right operand equals precedence of Or,
      // parenthesize s2, otherwise value remains s2. 
      val op2 = if (p2<=0) parens(s2) else s2;

      // Return tuple with first element the precedence of Or, and
      // second one the infix string representation of the Or-expression. 
      (0, op1+" Or "+op2)
    }


    def infix_And(l:(Int,String),r:(Int,String)) : (Int,String) = {
      // ...
      (0,"")


    }

    def infix_Not(o:(Int,String)) : (Int,String) = {
      // ...
      (0,"")


    }

    def infix_Const(v:Boolean) : (Int,String) = (0,"") // ...


    /* A minimal client of the traversal functionality */    

    val exp = Not(And(Or(Const(true),Const(false)),Or(Const(false),Const(true))));
    
    println ("eval Interpreter: "+eval(exp));

    println ("eval Visitor: "+exp.visit(eval_And,eval_Or,eval_Not,eval_Const));

    // ._2 retrieves second element of tuple
    println ("infix Visitor: "+exp.visit(infix_And,infix_Or,infix_Not,infix_Const)._2);

  }
}
 
  