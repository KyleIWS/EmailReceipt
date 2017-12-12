"use strict";

const CONNECT_URL = "http://localhost:4444"

$("#create").click(function() {
   requestGET("create").then(function(data, status, xhr) {
         console.log("I tried to create a new receipt");
         getAll();
      });
   });

$("#all").click(function() {
      getAll();
   });

$("#delete-all").click(function() {
   requestGET("delete-all").then(function(data, status, xhr) {
         console.log("deleting everything");
         getAll();
      });
   });

   function getAll() {
      requestGET("all").then(function(data, status, xhr) {
            $("#results").html("");
            let results = JSON.parse(data);
            console.log(results);
            for (var i = results.length - 1; i >= 0; i--) {
               let receiptData = results[i];
               let momentCreated = moment(receiptData.created);
               let receiptDiv = $("<div>", {"id": receiptData.id, "class": "row"});
               let imgTagText = $("<div>", {"class": "imgText"});
               imgTagText.text("<img src=" + CONNECT_URL + "/static/" + receiptData.id + ".png />");
               let createdDiv = $("<div>", {"class": "created"});
               createdDiv.html("Created on: " + momentCreated.format("MMM DD, YYYY") + " at " + momentCreated.format("hh:mm a"));
               let readsDiv = $("<div>", {"class": "reads"});
               for (var j = 0; j < receiptData.reads.length; j++) {
                     let momentRead = moment(receiptData.reads[j]);
                     let readDiv = $("<div>", {"class": "reads"});
                     readDiv.html("img tag was rendered " + momentRead.from(momentCreated, true) + " after the receipt was created. Time was " + momentRead.format("MMM DD, YYYY hh:mm a"));
                     readsDiv.append(readDiv);
               }
               receiptDiv.append(imgTagText).append(createdDiv).append(readsDiv);
               $("#results").append(receiptDiv);
            }
         });
   }

   function requestGET(path) {
      return $.ajax({
         url: CONNECT_URL + "/" + path,
         type: "GET",
      });
   }
