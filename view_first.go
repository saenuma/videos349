package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func drawFirstView(window *glfw.Window) {
	ProjObjCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &ProjObjCoords)

	theCtx.setFontSize(25)
	pnIRect := theCtx.drawInput(PROJ_NameInput, 30, 20, 420, "enter project name", false)
	pnBtnX := nextX(pnIRect, 30)
	nPRS := theCtx.drawButtonA(PROJ_NewProject, pnBtnX, 20, "New Project", fontColor, "#D5B5D2")
	oWDBX := nextX(nPRS, 140)
	wDBRS := theCtx.drawButtonA(PROJ_OpenWDBtn, oWDBX, 20, "Open Folder", "#fff", "#693E68")
	lS3X := nextX(wDBRS, 10)
	theCtx.drawButtonA(PROJ_LaunchS349, lS3X, 20, "V349 Slides", "#fff", "#693E68")
	theCtx.setFontSize(20)

	// second row border
	borderY := nextY(pnIRect, 30)

	fontPath := GetDefaultFontPath()
	theCtx.ggCtx.LoadFontFace(fontPath, 30)
	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawString("Continue Projects", 20, float64(borderY)+12+30)
	theCtx.ggCtx.LoadFontFace(fontPath, 20)

	projectFiles := GetProjectFiles()
	currentX := 40
	currentY := borderY + 22 + 30 + 10
	for i, pf := range projectFiles {

		btnId := 1000 + (i + 1)
		pfRect := theCtx.drawButtonA(btnId, currentX, currentY, pf.Name, "#fff", "#744A6F")

		newX := currentX + pfRect.Width + 10
		if newX > (wWidth - pfRect.Width) {
			currentY += 50
			currentX = 40
		} else {
			currentX += pfRect.Width + 10
		}

	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func fVMouseCB(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range ProjObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	rootPath, _ := GetRootPath()

	switch widgetCode {
	case PROJ_NewProject:
		if NameInputEnteredTxt == "" {
			return
		}

		// create file
		ProjectName = NameInputEnteredTxt + ".v3p"
		outPath := filepath.Join(rootPath, ProjectName)
		os.WriteFile(outPath, []byte(""), 0777)

		// move to work view
		drawItemsView(window, 1)
		window.SetMouseButtonCallback(iVMouseBtnCB)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(iVScrollBtnCB)
		// quick hover effect
		window.SetCursorPosCallback(getHoverCB(ObjCoords))

	case PROJ_OpenWDBtn:
		rootPath, _ := GetRootPath()
		ExternalLaunch(rootPath)

	case PROJ_LaunchS349:
		cmdPath := GetS349Command()
		cmd := exec.Command(cmdPath)
		cmd.Start()
	}

	if widgetCode > 1000 && widgetCode < 2000 {
		num := widgetCode - 1000 - 1
		projectFile := GetProjectFiles()[num]

		ProjectName = projectFile.Name

		// load instructions
		obj := make([]map[string]string, 0)
		rootPath, _ := GetRootPath()
		inPath := filepath.Join(rootPath, ProjectName)
		rawBytes, _ := os.ReadFile(inPath)
		json.Unmarshal(rawBytes, &obj)

		Instructions = append(Instructions, obj...)

		// move to work view
		drawItemsView(window, 1)
		window.SetMouseButtonCallback(iVMouseBtnCB)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(iVScrollBtnCB)
		window.SetCursorPosCallback(getHoverCB(ObjCoords))
	}
}

func fVKeyCB(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	wWidth, wHeight := window.GetSize()

	if key == glfw.KeyBackspace && len(NameInputEnteredTxt) != 0 {
		NameInputEnteredTxt = NameInputEnteredTxt[:len(NameInputEnteredTxt)-1]
	} else if key == glfw.KeySpace {
		NameInputEnteredTxt += " "
	} else if key == glfw.KeyEnter && len(NameInputEnteredTxt) != 0 {
		// create file
		rootPath, _ := GetRootPath()

		ProjectName = NameInputEnteredTxt + ".v3p"
		outPath := filepath.Join(rootPath, ProjectName)
		os.WriteFile(outPath, []byte(""), 0777)

		// move to work view
		drawItemsView(window, 1)
		window.SetMouseButtonCallback(iVMouseBtnCB)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(iVScrollBtnCB)
		return
	} else {
		NameInputEnteredTxt += glfw.GetKeyName(key, scancode)
	}

	nIRS := ProjObjCoords[PROJ_NameInput]
	theCtx := Continue2dCtx(CurrentWindowFrame, &ProjObjCoords)
	theCtx.setFontSize(25)
	theCtx.drawInput(PROJ_NameInput, nIRS.OriginX, nIRS.OriginY, nIRS.Width, NameInputEnteredTxt, true)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}
