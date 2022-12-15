package importmatcher

import (
	"fmt"
	"testing"
)

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
	fmt.Println("Found imports:")
	for _, foundImport := range foundImports {
		fmt.Println(foundImport)
	}
}
