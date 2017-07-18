package repositories

import (
	"net/http"
    "stufftobedone/domain"
    "io/ioutil"
    "log"
)

type CustomElementRepository struct {
  request    *http.Request
}

func NewCustomElementRepository(request *http.Request) *CustomElementRepository {
	r := new(CustomElementRepository)
	r.request = request
	return r
}

func (r *CustomElementRepository) GetByName(elementName string) (domain.CustomElement, error) {
    customElement := domain.CustomElement{}
    if elementName == "/x-element-to-do.html" {
        customElement.Body = `
            <dom-module id="element-to-do">
                <template>
                    <style>
                    :host {
                        display: block;
                    }
                    p {
                        background-color: orange;
                        color: white;
                        margin: 20px;
                    }
                    </style>
                    <p>Sample Element [[elementData.text]]</p>
                    <input type="text" value="{{elementData.text::input}}">
                    <paper-button on-tap="save">
                    <span>test</span>
                    </paper-button>
                </template>

                <script>
                    Polymer({

                    is: 'element-to-do',

                    properties: {
                        elementData: {
                            type: Object,
                            value: {text: 'hello', status: 'OPEN'},
                        },
                        isCompleted: {
                            type: Boolean,
                            value: false,
                        }
                    },
                    save: function(){
                        var result =  {data: this.elementData, isCompleted: this.isCompleted};
                        console.log('element-changed', result);
                        this.fire('element-changed', result, { bubbles: false });
                    },

                    });
                </script>
                </dom-module>
        `
    } else  if elementName == "/x-element-meeting.html" {
        customElement.Body = `
            <dom-module id="element-meeting">
                <template>
                    <style>
                    :host {
                        display: block;
                    }
                    p {
                        background-color: blue;
                        color: white;
                        margin: 20px;
                    }
                    p.meeting-line {
                        background-color: lightblue;
                        color: white;
                        margin: 5px;
                    }
                    </style>
                    <p>Meeting</p>
                    <template is="dom-repeat" items="[[elementData.items]]">
                        <p class="meeting-line">[[item.text]]</p>
                    </template>
                    <input type="text" placeholder="add meeting note">
                    <paper-button on-tap="save">
                    <span>test</span>
                    </paper-button>
                </template>

                <script>
                    Polymer({

                    is: 'element-meeting',

                    properties: {
                        elementData: {
                            type: Object,
                            value: {items: []},
                        },
                        isCompleted: {
                            type: Boolean,
                            value: false
                        }
                    },
                    save: function(){
                        var result =  {data: this.elementData, isCompleted: this.isCompleted};
                        console.log('element-changed', result);
                        this.fire('element-changed', result, { bubbles: false });
                    },

                    });
                </script>
                </dom-module>
        `
    } else  { //} if elementName == "/element-test.html" {
        // read the element from the file
        buf, err := ioutil.ReadFile("./public/elements/default" + elementName)
        if err != nil {
            log.Println("read element failed", err) 
        }
        customElement.Body = string(buf)
    }
	return customElement, nil
}

func (r *CustomElementRepository) FindBookElements(bookID string) ([]domain.BookElement, error) {
    toDoElement := domain.BookElement{ Name: "To-do", ElementName: "element-to-do", Icon: "check-box" }
    noteElement := domain.BookElement{ Name: "Note", ElementName: "element-note", Icon: "speaker-notes" }
    meetingElement := domain.BookElement{ Name: "Meeting", ElementName: "element-meeting", Icon: "social:people" }
    sketchElement := domain.BookElement{ Name: "Sketch", ElementName: "element-sketch", Icon: "gesture" }
    //testElement := domain.BookElement{ Name: "Test", ElementName: "element-test", Icon: "bug-report" }
    
    bookElements := []domain.BookElement{ toDoElement, noteElement, meetingElement, sketchElement } //, testElement}
    
    
    return bookElements, nil
}