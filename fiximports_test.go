package autoimport

import (
	"testing"
)

// Mock implementation of ImportMatcher for testing
type MockImportMatcher struct{}

func (m *MockImportMatcher) StarPath(word string) (string, string) {
	switch word {
	case "ArrayList":
		return "ArrayList", "java.util.*"
	case "Map":
		return "Map", "java.util.*"
	default:
		return "", ""
	}
}

func TestFixImports(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       string
		removeExisting bool
	}{
		{
			name: "simple Java imports",
			input: `
package com.example;

import java.util.Map;
import java.util.ArrayList;

public class Main {
    ArrayList<String> list = new ArrayList<>();
    Map<String, String> map;
}
`,
			expected: `
package com.example;

import java.util.*; // ArrayList, Map

public class Main {
    ArrayList<String> list = new ArrayList<>();
    Map<String, String> map;
}
`,
			removeExisting: true,
		},
		{
			name: "gl",
			input: `
import net.java.games.jogl.*;

import java.awt.*;
import java.awt.event.*;

import javax.swing.*;

/** This is a basic JOGL app. Feel free to reuse this code or modify it. */
public class SimpleJoglApp extends JFrame {
    public static void main(String[] args) {
        final SimpleJoglApp app = new SimpleJoglApp();

        // show what we've done
        SwingUtilities.invokeLater(
                new Runnable() {
                    public void run() {
                        app.setVisible(true);
                    }
                });
    }

    public SimpleJoglApp() {
        // set the JFrame title
        super("Simple JOGL Application");

        // kill the process when the JFrame is closed
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

        // only three JOGL lines of code ... and here they are
        GLCapabilities glcaps = new GLCapabilities();
        GLCanvas glcanvas = GLDrawableFactory.getFactory().createGLCanvas(glcaps);
        glcanvas.addGLEventListener(new SimpleGLEventListener());

        // add the GLCanvas just like we would any Component
        getContentPane().add(glcanvas, BorderLayout.CENTER);
        setSize(500, 300);

        // center the JFrame on the screen
        centerWindow(this);
    }

    public void centerWindow(Component frame) {
        Dimension screenSize = Toolkit.getDefaultToolkit().getScreenSize();
        Dimension frameSize = frame.getSize();

        if (frameSize.width > screenSize.width) frameSize.width = screenSize.width;
        if (frameSize.height > screenSize.height) frameSize.height = screenSize.height;

        frame.setLocation(
                (screenSize.width - frameSize.width) >> 1,
                (screenSize.height - frameSize.height) >> 1);
    }
}
`,
			expected: `
import java.awt.*; // Dimension
import java.awt.event.*;
import javax.swing.*;
import net.java.games.jogl.*;


/** This is a basic JOGL app. Feel free to reuse this code or modify it. */
public class SimpleJoglApp extends JFrame {
    public static void main(String[] args) {
        final SimpleJoglApp app = new SimpleJoglApp();

        // show what we've done
        SwingUtilities.invokeLater(
                new Runnable() {
                    public void run() {
                        app.setVisible(true);
                    }
                });
    }

    public SimpleJoglApp() {
        // set the JFrame title
        super("Simple JOGL Application");

        // kill the process when the JFrame is closed
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

        // only three JOGL lines of code ... and here they are
        GLCapabilities glcaps = new GLCapabilities();
        GLCanvas glcanvas = GLDrawableFactory.getFactory().createGLCanvas(glcaps);
        glcanvas.addGLEventListener(new SimpleGLEventListener());

        // add the GLCanvas just like we would any Component
        getContentPane().add(glcanvas, BorderLayout.CENTER);
        setSize(500, 300);

        // center the JFrame on the screen
        centerWindow(this);
    }

    public void centerWindow(Component frame) {
        Dimension screenSize = Toolkit.getDefaultToolkit().getScreenSize();
        Dimension frameSize = frame.getSize();

        if (frameSize.width > screenSize.width) frameSize.width = screenSize.width;
        if (frameSize.height > screenSize.height) frameSize.height = screenSize.height;

        frame.setLocation(
                (screenSize.width - frameSize.width) >> 1,
                (screenSize.height - frameSize.height) >> 1);
    }
}
`,
			removeExisting: false,
		},
	}

	const onlyJava = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ima, err := New(onlyJava, tt.removeExisting)
			if err != nil {
				t.Fatalf("Error initializing the import matcher: %v", err)
			}
			const verbose = true
			got, err := ima.FixImports([]byte(tt.input), verbose)
			if err != nil {
				t.Fatalf("Error processing FixImports: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, got)
				//os.WriteFile("expected.java", []byte(tt.expected), 0644)
				//os.WriteFile("got.java", []byte(got), 0644)
			}
		})
	}
}
