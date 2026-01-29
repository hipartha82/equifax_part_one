// Golang program online for dashboard - file based input processing

package main

import "fmt"
import "os"
import "bufio"
import "strings"
import "strconv"
import "context"
import "time"
import "encoding/json"
import "sync"

type UserDetailsAndTodo struct
  {
    Id string `json:"id"`
    Full_name string `json:"full_name"`
    Status string `json:"status"`
    Pending_task_count string `json:"pending_task_count"`
    Next_urgent_task string `json:"next_urgent_task"`
    Error_warning string `json:"error_warning"`
  }
type userResult struct {
    userdata  []string
    founduser bool
    errfromfunc error
  }

func main() {
 
 var userinfo [][]string
 var todoinfo [][]string
 
 file, err := os.Create("cust_data.txt")
 if (err != nil) {
    println ("File cannot be created")
    return
 } else {
    println ("Details File Created")
    file.WriteString("1,John, Paul, 40 \n2,Sameer, kumar, 56 \n3,Rahul, Paul,  49\n4 , Raymond, Lee, 58 \n5 , Rakesh, Chauhan, 38 ")
    defer file.Close()
 }  
 
 file2, err2 := os.Create("cust_todo.txt")
 if (err2 != nil) {
    println ("File cannot be created")
    return
 } else {
    println ("ToDo File Created")
    file2.WriteString("1,HLD to Write, Y \n1,Dev Process to build, Y \n1,Testing to certify, N \n 2, HLD to Write, N \n2,Dev Process to build, N \n2,Testing To certify, N \n 3, HLD to Write, Y \n3,Dev Process to build, N \n3,Testing To certify, N \n 4, HLD to Write, Y \n4,Dev Process to build, Y \n4,Testing To certify, Y \n6, HLD to Write, N \n6,Dev Process to build, N \n6,Testing To certify, Y")
    defer file2.Close()
 }  

 userdetails, errdata := os.Open("cust_data.txt")
 if (errdata != nil) {
    println ("File not present or cannot be read")
    return
 } else {
     userinfo = readAndStoreUser(userdetails)
     defer userdetails.Close()
 
    println ("\nExpected output of User data from the array:")
    println ("------------------------------------------------")
    for i := 0; i < len(userinfo) ; i++ {
        for j := 0 ; j < len(userinfo[i]) ; j++ {
            print (userinfo[i][j],"#")
        }
        println ("")
    }
  }
  
  usertodo, errtodo := os.Open("cust_todo.txt")
  if (errtodo != nil) {
    println ("File not present or cannot be read")
    return
  } else {
     todoinfo = readAndStoreTodo(usertodo)
     defer usertodo.Close()
 
    println ("\nExpected output of Todo data from the array:")
    println ("------------------------------------------------")
    for i := 0; i < len(todoinfo) ; i++ {
        for j := 0 ; j < len(todoinfo[i]) ; j++ {
            print (todoinfo[i][j],"#")
        }
        println ("")
    }
  }

  var inpUserID string
  var udt []string
  var tds []string
  var foundudt, foundtds bool
  var errudt, errtds error
  var result UserDetailsAndTodo
  var udetails, utodos userResult
  var waitGrp sync.WaitGroup
  
  inpUserID = "3"
  println ("\nThe User ID to get the data for details and Todo is : ", inpUserID)
  
  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  
  waitGrp.Add(2)

  go func() {
      defer waitGrp.Done();
      udt, foundudt, errudt = findDetailsForUser(ctx,inpUserID,userinfo)
  }()
  
  go func() {
    defer waitGrp.Done();
    tds, foundtds, errtds = findTodoForUser(ctx,inpUserID,todoinfo)
  }()
  
  waitGrp.Wait()
  
  udetails = userResult{udt, foundudt, errudt}
  utodos = userResult{tds, foundtds, errtds}
  
  result = formatOutputToJson(udetails,utodos)
  jsonOutput, _ := json.MarshalIndent(result, "", "  ")
  println ("\nThe Merged output for the User details with Todo Summary is : \n", string(jsonOutput))
 
} //end of main


func readAndStoreUser (f *os.File) ([][]string) {
  var uinfo [][]string
  scanner := bufio.NewScanner(f)
  println ("\nOriginal User file entries:")
  println ("-------------------------------")
  for scanner.Scan() {
    fmt.Print(scanner.Text(),"  ")
    fullname := ""
    var temp []string 
    
    userdata := strings.Split(scanner.Text(),"," )
    for i, v := range userdata {
        if (i == 1 ) {
            fullname += strings.TrimSpace(v) + " "
            } else if (i == 2) {
                fullname += strings.TrimSpace(v)
            } else if (i == 3) {
                age,_ := strconv.Atoi(strings.TrimSpace(v)) 
                if (age > 50) {
                    temp = append(temp,"Veteran")
                    //print ("Veteran")
                } else {
                    temp = append(temp,"Rookie")
                    //print ("Rookie")
                }
            }
        
        switch i {
            case 0 :
              //print (strings.TrimSpace(v), ",")
              temp = append(temp,strings.TrimSpace(v))
            case 2 :
              //print (fullname, ",")
              temp = append(temp,fullname)
        }
    }
    println("")
    uinfo = append(uinfo,temp)
 }
 return uinfo
}

func readAndStoreTodo (f *os.File) ([][]string) {
  var tdinfo [][]string
  scanner := bufio.NewScanner(f)
  println ("\nOriginal Todo file entries:")
  println ("-------------------------------")
  
  for scanner.Scan() {
    fmt.Print(scanner.Text(),"  ")
    var temp []string 
    sameuser := false
    var userindex int
    var tdesc string
    
    tddata := strings.Split(scanner.Text(),"," )
    for i, v := range tddata {
        // identify if the current user id is already present in tdinfo
        if len(tdinfo) > 0 && i == 0 {
            for itr := 0; itr < len(tdinfo) && !sameuser; itr++ {
                if (tdinfo[itr][0] == strings.TrimSpace(v)) {
                    sameuser = true
                    userindex = itr
                }
            }
        }
        // For same user update the existing entry in tdinfo else add new row in tdinfo from temp
        if sameuser { 
            if i == 1 {
                tdesc = strings.TrimSpace(v) // to be used for later if needed to update
            }
            if i == 2 && strings.TrimSpace(v) == "N" {
                x,_ := strconv.Atoi(tdinfo[userindex][2]) 
                x += 1
                tdinfo[userindex][2] = strconv.Itoa(x) 
                if tdinfo[userindex][2] == "1" {
                    tdinfo[userindex][1] = tdesc
                }
            }
        } else {
            if i == 2 && strings.TrimSpace(v) == "N" {
                temp = append(temp,"1") 
            } else if i == 2 && strings.TrimSpace(v) == "Y" {
                temp = append(temp,"0")
                temp[i-1] = ""
            } else {
                temp = append(temp,strings.TrimSpace(v))
            }
        }
      } //processing of one input line from the Todo file is finished
     
    println("")
    if !sameuser {
        tdinfo = append(tdinfo,temp)   
    }
  }
  return tdinfo
 }
 
 func findDetailsForUser (cntx context.Context, userId string, userDet [][]string) ([]string,bool,error) {
    userId = strings.TrimSpace(userId)
    for i := 0; i < len(userDet) ; i++ {
        select {
            case <- cntx.Done() :
                return nil, false, cntx.Err()
            default :
                if userDet[i][0] == userId {
                    /*
                    print ("\n\nThe matched user details are : ")
                    for j := 0 ; j < len(userDet[i]) ; j++ {
                        print (userDet[i][j], " # " )
                    } 
                    */
                    return userDet[i], true, nil
                }
        }
    }
    return nil, false, nil
 }
 
  func findTodoForUser (cntx context.Context, userId string, todoSummary [][]string) ([]string,bool,error) {
    userId =  strings.TrimSpace(userId)
    for i := 0; i < len(todoSummary) ; i++ {
        select {
            case <- cntx.Done() :
                 return nil, false, cntx.Err()
            default :
                if todoSummary[i][0] == userId {
                    /*
                    print ("\n\nToDo Summary for the matched user are : ")
                    for j := 0 ; j < len(todoSummary[i]) ; j++ {
                        print (todoSummary[i][j], " # " )
                    } 
                    */
                    return todoSummary[i],true,nil
                }
        }
    }
    return nil, false, nil
 }
 
 func formatOutputToJson (usrdtls userResult, usrtodos userResult) UserDetailsAndTodo {
     
     var jsnop UserDetailsAndTodo
     
     if (usrdtls.founduser && usrtodos.founduser) {
      println ("\n\nBoth user details and Todo Summary are successfully received")
      fmt.Println("User Details : ",usrdtls)
      fmt.Println("Todo Summary : ",usrtodos)
      jsnop = UserDetailsAndTodo { 
          Id : usrdtls.userdata[0],
          Full_name : usrdtls.userdata[1],
          Status : usrdtls.userdata[2],
          Pending_task_count : usrtodos.userdata[2],
          Next_urgent_task : usrtodos.userdata[1],
          Error_warning : "null",
        }
      } else if (usrdtls.founduser && !usrtodos.founduser) {
          if usrtodos.errfromfunc == context.DeadlineExceeded {
              println ("\n\nTransaction timed out while fetching ToDo Info for the user")
              fmt.Println("User Details : ",usrdtls)
              jsnop = UserDetailsAndTodo {
                  Id : usrdtls.userdata[0],
                  Full_name : usrdtls.userdata[1],
                  Status : usrdtls.userdata[2],
                  Pending_task_count : "null",
                  Next_urgent_task : "null",
                  Error_warning : "ToDo unavailable due to ToDo transaction got timed out",
                }
          } else {
              println ("\n\nOnly User details is successfully received, Todo Info is missing in the source data")
              fmt.Println("User Details : ",usrdtls)
              jsnop = UserDetailsAndTodo {
                  Id : usrdtls.userdata[0],
                  Full_name : usrdtls.userdata[1],
                  Status : usrdtls.userdata[2],
                  Pending_task_count : "null",
                  Next_urgent_task : "null",
                  Error_warning : "User ID is not available in the source Todo dataset",
                }
          }
      } else if (!usrdtls.founduser && usrtodos.founduser) {
          if usrdtls.errfromfunc == context.DeadlineExceeded {
              println ("\n\nTransaction timed out while fetching User details Info for the user")
              fmt.Println("ToDo Summary : ",usrtodos)
              jsnop = UserDetailsAndTodo{
                  Id : "null",
                  Full_name : "null",
                  Status : "null",
                  Pending_task_count : usrtodos.userdata[2],
                  Next_urgent_task : usrtodos.userdata[1],
                  Error_warning : "User details unavailable due to user details transaction got timed out",
                }
          } else {
              println ("\n\nOnly ToDo details is successfully received, User Details info is missing in the source data")
              fmt.Println("ToDo Summary : ",usrtodos)
              jsnop = UserDetailsAndTodo{
                  Id : "null",
                  Full_name : "null",
                  Status : "null",
                  Pending_task_count : usrtodos.userdata[2],
                  Next_urgent_task : usrtodos.userdata[1],
                  Error_warning : "User ID is not available in the source User Details dataset",
                }
          }
      } else {
          if usrdtls.errfromfunc == context.DeadlineExceeded && usrtodos.errfromfunc == context.DeadlineExceeded {
              println ("\n\nWhole Transaction timed out while fetching both User details info and ToDo info for the user")
          } else {
              println ("\n\nNone of User ID and Todo details are available in the dataset for the user")
          }
      }
    return jsnop
 }
