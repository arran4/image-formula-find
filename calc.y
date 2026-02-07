%{
package image_formula_find;

import __yyfmt__ "fmt"

var yyResult *Function
%}

%token<expr> Highest
%token<float> FLOAT
%token<s> VAR FUNCNAME
%type<expr> expr

%union {
    float float64
    s string
    expr Expression
 }

%right '='
%left '+' '-'
%left '*' '/' '%' '^' ','
%right Highest FUNCNAME

%%
input
    : expr '=' expr { yyResult = &Function{ Equals: &Equals { LHS: $1, RHS: $3 } } }
    ;

expr: FLOAT             { $$ = &Const{Value: $1} }
    | VAR               { $$ = &Var{ Var: $1 } }
    | expr '+' expr     { $$ = &Plus{ LHS: $1, RHS: $3, } }
    | expr '-' expr     { $$ = &Subtract{ LHS: $1, RHS: $3, } }
    | expr '*' expr     { $$ = &Multiply{ LHS: $1, RHS: $3, } }
    | expr '/' expr     { $$ = &Divide{ LHS: $1, RHS: $3, } }
    | expr '%' expr     { $$ = &Modulus{ LHS: $1, RHS: $3, } }
    | expr '^' expr     { $$ = &Power{ LHS: $1, RHS: $3, } }
    | '+' expr  %prec Highest    { $$ = $2 }
    | '-' expr  %prec Highest    { $$ = &Negate{ Expr: $2 } }
    | expr FUNCNAME expr     { $$ = NewDoubleFunction($2, $1, $3, true) }
    | FUNCNAME '(' expr ')'  { $$ = NewSingleFunction($1, $3) }
    | FUNCNAME '(' expr ',' expr ')' { $$ = NewDoubleFunction($1, $3, $5, false) }
    | '(' expr ')'            { $$ = &Brackets{ Expr: $2 } }
    ;

%%
