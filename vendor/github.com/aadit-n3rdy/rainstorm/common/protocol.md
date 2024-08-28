# Rainstorm Protocol
 
 Rainstorm is a P2P file transfer protocol, consisting of trackers and peers. 
 All communication occurs through JSON.

## General Notes
 1. Each packet has a "class" and a "type"
 2. Packet classes are `init`, `tracking` and `ft`

## Init Packets

 1. `download_start`:
    - peer -> tracker
    - To retrieve the set of peers for a file
    - **Mandatory fields**:
        - `file_id`
    - **Responses**
        - `download_data`

2. `download_data`:
    - tracker -> peer
    - Contains the set of peers, and the checksum values
    - **Mandatory fields**:
        - `status`: 0 for success, 1 for failure
    - **Success fields**:
        - `data`: a JSON object which can be converted into a FileDownloadData struct
    - **Error fields**:
        - `error`: error message in case of failure
