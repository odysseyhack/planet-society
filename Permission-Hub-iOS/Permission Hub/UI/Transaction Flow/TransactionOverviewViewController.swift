//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionOverviewViewController: UITableViewController {

    // MARK: - Private properties

    private let transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        configureTableView()
    }

    // MARK: - Configuration

    private func configureTableView() {

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.allowsSelection = false
    }

    // MARK: - UITableViewDataSource

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return items.count
    }

    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        let cell = tableView.dequeueReusableCell(
            withIdentifier: String(describing: TransactionTableViewCell.self),
            for: indexPath)

        return cell
    }
}
