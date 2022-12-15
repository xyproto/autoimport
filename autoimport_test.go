package importmatcher

import (
	"testing"
)

const organizedImportsShouldLookLike = `import java.io.File;
import java.io.FileNotFoundException;
import java.util.Scanner;
`

const sourceCode = `
public class ReadFile {
  public static void main(String[] args) {
    try {
      File myObj = new File("filename.txt");
      Scanner myReader = new Scanner(myObj);
      while (myReader.hasNextLine()) {
        String data = myReader.nextLine();
        System.out.println(data);
      }
      myReader.close();
    } catch (FileNotFoundException e) {
      System.out.println("An error occurred.");
      e.printStackTrace();
    }
  }
}
`

func TestFindImports(t *testing.T) {
	impM, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundImports := impM.FindImports(sourceCode)
	if !hasS(foundImports, "java.io.File") {
		t.Errorf("The list of found imports should include java.io.File\n")
	}
	if !hasS(foundImports, "java.io.FileNotFoundException") {
		t.Errorf("The list of found imports should include java.io.FileNotFoundException\n")
	}
	if !hasS(foundImports, "java.util.Scanner") {
		t.Errorf("The list of found imports should include java.util.Scanner\n")
	}
	if len(foundImports) != 3 {
		t.Errorf("There should only be 3 found imports\n")
	}
}

func TestOrganizedImports(t *testing.T) {
	impM, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	organizedImports := impM.OrganizedImports(sourceCode, true)
	if organizedImports != organizedImportsShouldLookLike {
		t.Errorf("The organized imports looks like:\n\n%s\nBut they should look like:\n\n%s\n", organizedImports, organizedImportsShouldLookLike)
	}
}
